import * as path from 'node:path'
import { beforeEach, describe, expect, it } from 'vitest'
import type { Component } from '../../src/imports/fogg-types.generated'
import { loadComponentConfig } from '../../src/util/load-component-config'

let componentConfig: Component
describe('loadComponentConfig', () => {
  beforeEach(() => {
    componentConfig = loadComponentConfig(
      path.join(__dirname, 'fixtures', '.fogg-component.yaml'),
    )
  })

  it('loads the component configuration', () => {
    expect(componentConfig).toBeDefined()
  })

  it('returns an object with the expected properties', () => {
    expect(componentConfig).toMatchObject({
      component_backends: expect.any(Object),
      component_backends_filtered: expect.any(Boolean),
      required_providers: expect.any(Object),
    })
  })
})
