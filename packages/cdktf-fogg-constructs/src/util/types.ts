/*
 * Helper type to make readonly interface properties mutable
 */
export type Mutable<T> = {
  -readonly [P in keyof T]: T[P]
}

/**
 * Helper type to make all properties of an interface of type string
 */
export type StringifiedRecord<T> = { [K in keyof T]: string }
