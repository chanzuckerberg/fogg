import * as path from 'node:path'
import { TerraformOutput, Testing } from 'cdktf'
import { describe, expect, it } from 'vitest'
import {
  SchemaVitestMatchers,
  setupSchemaMatchers,
} from '../../src/matchers/schema'

// satisfy TypeScript :(
declare module 'vitest' {
  interface Assertion<T = any> extends SchemaVitestMatchers<T> {}
}

await setupSchemaMatchers()
describe('Output Schema matching', () => {
  const sourceFilePath = path.resolve(__dirname, './fixtures/output-schema.ts')
  it('extends imported `expect` with Schema matchers', () => {
    const assertion = expect(true)
    expect(assertion.toMatchSchema).toBeTypeOf('function')
  })

  it('passes if outputs matches TypeScript interface', () => {
    const output = { foo: 'aws_example.foo.arn', bar: 'aws_example.bar.arn' }
    expect(() => expect(output).toMatchSchema(sourceFilePath, 'OutputSchema'))
  })

  it('fails if outputs does not match TypeScript interface (missing keys)', () => {
    const output = { foo: 'aws_example.foo.arn' }
    expect(() =>
      expect(output).toMatchSchema(sourceFilePath, 'OutputSchema'),
    ).toThrowError(/Expected object to match schema/)
  })
  it('fails if outputs does not match TypeScript interface (extra keys)', () => {
    const output = {
      foo: 'aws_example.foo.arn',
      bar: 'aws_example.bar.arn',
      dynamic: 'aws_example.bar.arn',
    }
    expect(() =>
      expect(output).toMatchSchema(sourceFilePath, 'OutputSchema'),
    ).toThrowError(/Expected object to match schema/)
  })

  it('passes if outputs matches TypeScript interface with ignored keys', () => {
    const output = {
      foo: 'aws_example.foo.arn',
      bar: 'aws_example.bar.arn',
      dynamic: 'aws_example.bar.arn',
    }
    const ignoredKeys = ['dynamic']
    expect(() =>
      expect(output).toMatchSchema(sourceFilePath, 'OutputSchema', ignoredKeys),
    )
  })

  it('passes with stack outputs', () => {
    const synthed = Testing.synthScope((stack) => {
      new TerraformOutput(stack, 'foo', { value: 'foo' })
      new TerraformOutput(stack, 'bar', { value: 'foo' })
    })
    const outputs = JSON.parse(synthed).output
    expect(outputs).toBeDefined
    expect(outputs).toMatchSchema(sourceFilePath, 'OutputSchema')
  })
})
