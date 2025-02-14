import * as fs from 'node:fs'
import * as yaml from 'js-yaml'
import { Component } from '../imports/fogg-types.generated'

export function loadComponentConfig(
  path: string = `.fogg-component.yaml`,
): Component {
  const file = fs.readFileSync(path, 'utf8')
  const componentConfig = yaml.load(file) as Component
  return replaceNullWithUndefined(componentConfig)
}

// helper function to replace fogg "null" values with "undefined"
function replaceNullWithUndefined(obj: any): any {
  if (obj === null) {
    return undefined
  }
  if (Array.isArray(obj)) {
    return obj.map(replaceNullWithUndefined)
  }
  if (typeof obj === 'object' && obj !== null) {
    const newObj: any = {}
    for (const key in obj) {
      newObj[key] = replaceNullWithUndefined(obj[key])
    }
    return newObj
  }
  return obj
}
