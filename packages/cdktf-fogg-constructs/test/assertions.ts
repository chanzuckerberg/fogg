import {
  AnnotationMetadataEntryType,
  StackAnnotation,
  TerraformStack,
  Testing,
} from 'cdktf'
import { TerraformConstructor } from 'cdktf/lib/testing/matchers'
import { MetadataEntry } from 'constructs'
import { Assertion, expect } from 'vitest'

export interface SynthOptions {
  /**
   * snapshot full synthesized template
   */
  snapshot?: boolean
  /**
   * Run all validations on the stack before synth
   */
  runValidations?: boolean
}

/**
 * Helper class to create Jest Matchers for a TerraformStack
 */
export class Template {
  /**
   * Create Vitest Assertions from the parsed synthesized spec
   */
  static fromStack(
    stack: TerraformStack,
    options: SynthOptions = {},
  ): Assertion<any> {
    const synthesized = Template.getSynthString(stack, options)
    const parsed = JSON.parse(synthesized)
    return expect(parsed)
  }
  /**
   * Create Vitest Assertions for the synthesized JSON string
   *
   * This always runs TerraformStack.prepareStack() as this
   * library heavily depends on it for pre-synth resource
   * generation.
   */
  static synth(
    stack: TerraformStack,
    options: SynthOptions = {},
  ): Assertion<any> {
    const synthesized = Template.getSynthString(stack, options)
    return expect(synthesized)
  }

  /**
   * Create Vitest Assertions for stack resources of a specific type
   *
   * This always runs TerraformStack.prepareStack() as this
   * library heavily depends on it for pre-synth resource
   * generation.
   */
  static resources(
    stack: TerraformStack,
    type: TerraformConstructor,
    options: SynthOptions = {},
  ): Assertion<any> {
    const synthesized = Template.getSynthString(stack, options)
    const parsed = JSON.parse(synthesized)
    const resources = parsed.resource
      ? Object.values(parsed.resource[type.tfResourceType] ?? {})
      : []
    return expect(resources)
  }

  /**
   * Create Vitest Assertions for stack resources of a specific type
   *
   * This always runs TerraformStack.prepareStack() as this
   * library heavily depends on it for pre-synth resource
   * generation.
   */
  static dataSources(
    stack: TerraformStack,
    type: TerraformConstructor,
    options: SynthOptions = {},
  ): Assertion<any> {
    const synthesized = Template.getSynthString(stack, options)
    const parsed = JSON.parse(synthesized)
    const dataSources = parsed.data
      ? Object.values(parsed.data[type.tfResourceType] ?? {})
      : []
    return expect(dataSources)
  }

  private static getSynthString(
    stack: TerraformStack,
    options: SynthOptions = {},
  ): string {
    const { snapshot = false, runValidations = false } = options
    stack.prepareStack() // required to add pre-synth resources
    const result = Testing.synth(stack, runValidations)
    if (snapshot) {
      expect(result).toMatchSnapshot()
    }
    return result
  }
}

export class Annotations {
  public static fromStack(stack: TerraformStack): Annotations {
    // https://github.com/hashicorp/terraform-cdk/blob/v0.20.10/packages/cdktf/lib/synthesize/synthesizer.ts#L59-L74
    // collect Annotations into Manifest
    const annotations = stack.node
      .findAll()
      .map((node) => ({
        node,
        metadatas: node.node.metadata.filter(isAnnotationMetadata),
      }))
      .map<StackAnnotation[]>(({ node, metadatas }) =>
        metadatas.map((metadata) => ({
          constructPath: node.node.path,
          level: metadata.type as AnnotationMetadataEntryType,
          message: metadata.data,
          stacktrace: metadata.trace,
        })),
      )
      .reduce((list, metadatas) => [...list, ...metadatas], []) // Array.flat()
    return new Annotations(annotations)
  }

  private constructor(private readonly annotations: StackAnnotation[]) {}

  public get warnings(): StackAnnotation[] {
    return this.annotations.filter(isWarningAnnotation)
  }
  public get errors(): StackAnnotation[] {
    return this.annotations.filter(isErrorAnnotation)
  }

  /**
   * check if the stack has a warning for certain context path and message
   */
  public hasWarnings(
    ...expectedWarnings: Array<Partial<StackAnnotationMatcher>>
  ) {
    const warningMatchers = expectedWarnings.map((warning) => {
      const transformed: Partial<StackAnnotationMatcher> = {}
      for (const key in warning) {
        const value = warning[key as keyof StackAnnotationMatcher]
        if (value instanceof RegExp) {
          transformed[key as keyof StackAnnotationMatcher] =
            expect.stringMatching(value)
        } else {
          transformed[key as keyof StackAnnotationMatcher] = value
        }
      }
      return expect.objectContaining(transformed)
    })
    expect(this.warnings).toEqual(expect.arrayContaining(warningMatchers))
  }
}

// https://github.com/hashicorp/terraform-cdk/blob/v0.20.10/packages/cdktf/lib/synthesize/synthesizer.ts#L164-L173
const annotationMetadataEntryTypes = [
  AnnotationMetadataEntryType.INFO,
  AnnotationMetadataEntryType.WARN,
  AnnotationMetadataEntryType.ERROR,
] as string[]
function isAnnotationMetadata(metadata: MetadataEntry): boolean {
  return annotationMetadataEntryTypes.includes(metadata.type)
}

function isErrorAnnotation(annotation: StackAnnotation): boolean {
  return annotation.level === AnnotationMetadataEntryType.ERROR
}

function isWarningAnnotation(annotation: StackAnnotation): boolean {
  return annotation.level === AnnotationMetadataEntryType.WARN
}

export interface StackAnnotationMatcher {
  constructPath: string | RegExp
  message: string | RegExp
}
