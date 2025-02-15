import { DataTerraformRemoteState, StringListMap, Token } from 'cdktf'
import { Fn } from 'cdktf'

export interface ListOutput<T = string> {
  /**
   * Returns a token that resolves to the element at `index`.
   */
  element(index: number): T

  /**
   * Returns the remote state token, use Fn.element to access elements.
   */
  readonly token: string[]
}

/**
 * A small wrapper for map outputs.
 *
 * T: interface representing the keys
 * V: the type of each map value (often string, but can be something else).
 */
export interface MapOutput<T, K extends keyof T = keyof T, V = string> {
  /**
   * Looks up the given key in the map, returning a Token
   * that resolves to the value.
   */
  lookup(key: K, defaultValue?: V): V
  /**
   * returns a property access expression that accesses the property at the given path.
   * For example lookupNested("a", "b", "c") will return a Terraform expression like x["a"]["b"]["c"]
   */
  lookupNested(...path: string[]): V

  /**
   * Returns the entire map as a Token.
   */
  readonly token: string
}

/**
 * Infer output schema from a given type T against user provided schema.
 */
export type InferOutputSchema<T> = {
  [K in keyof T]: T[K] extends ListOutput<any>
    ? 'list'
    : T[K] extends MapOutput<any, any>
      ? 'map'
      : T[K] extends string
        ? 'string'
        : T[K] extends number
          ? 'number'
          : T[K] extends boolean
            ? 'boolean'
            : never // Or fallback if you have other shapes
}

export class RemoteStateAccessProxy<T extends Record<string, any>> {
  private readonly proxy: T

  constructor(
    private readonly data: DataTerraformRemoteState,
    private readonly schema?: InferOutputSchema<T>,
  ) {
    this.proxy = new Proxy({} as T, {
      get: (_target, prop: string | symbol) => {
        if (typeof prop !== 'string') {
          return undefined // ignore symbol or non-string props
        }

        // If no schema is provided, default to getString:
        if (!this.schema) {
          return this.data.getString(prop)
        }

        const declaredType = this.schema[prop as keyof T]
        switch (declaredType) {
          case 'string':
            return this.data.getString(prop)
          case 'list': {
            const listWrapper: ListOutput<T> = {
              element: (index: number) => {
                return Fn.element(this.data.getList(prop), index)
              },
              token: this.data.getList(prop),
            }
            return listWrapper
          }
          // TODO: nested proxy to throw on incorrectly accessing list items
          // to validate numeric indexing, return a proxy which throws on numeric props
          // if (/^\d+$/.test(prop)) {
          //   throw new Error(
          //     'Numeric indexing remote state outputs is not supported. Either use Fn.element(..) or provide outputSchema as "list" and use element(...) instead.'
          //   );
          // }
          case 'map': {
            const mapWrapper: MapOutput<T> = {
              lookup: (key: string, defaultValue?: string) =>
                Fn.lookup(this.data.get(prop), key, defaultValue ?? ''),
              lookupNested: (...path: string[]) =>
                Fn.lookupNested(this.data.get(prop), path),
              token: this.data.getString(prop),
            }
            return mapWrapper
          }
          case 'number':
            return this.data.getNumber(prop)
          case 'boolean':
            return this.data.getBoolean(prop)
          default:
            return this.data.get(prop)
        }
      },
    })
  }

  /**
   * Expose the proxied object, which “looks like” T
   * but each property is actually a token from the remote state.
   */
  public asObject(): T {
    return this.proxy
  }
}
