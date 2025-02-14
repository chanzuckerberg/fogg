import {
  ModuleKind,
  ModuleResolutionKind,
  Project,
  ProjectOptions,
  ScriptTarget,
} from 'ts-morph'

// Create a minimal compiler options configuration
const compilerOptions: ProjectOptions = {
  compilerOptions: {
    strict: true,
    target: ScriptTarget.ES2022,
    module: ModuleKind.Node16,
    moduleResolution: ModuleResolutionKind.Node16,
  },
  skipAddingFilesFromTsConfig: true, // Skip loading files from tsconfig
  skipFileDependencyResolution: true, // Skip resolving file dependencies
}

// Initialize project once
const project = new Project(compilerOptions)

/**
 * Vitest Matcher for a TypeScript interface against an object
 *
 * @param received
 * @param sourceFilePath
 * @param interfaceName
 * @param ignoredKeys
 * @returns
 */
export function toMatchSchema(
  received: any,
  sourceFilePath: string,
  interfaceName: string,
  ignoredKeys: string[] = [],
) {
  const sourceFile =
    project.getSourceFile(sourceFilePath) ??
    project.addSourceFileAtPath(sourceFilePath)
  if (!sourceFile) {
    throw new Error('Could not find source file')
  }

  if (!sourceFile) {
    throw new Error('Could not find source file')
  }

  const interfaceDecl = sourceFile.getInterfaceOrThrow(interfaceName)
  const schemaKeySet = new Set(
    interfaceDecl.getProperties().map((p) => {
      const nameNode = p.getNameNode()
      // If the property name is a string literal, return its literal value without quotes.
      if (nameNode.getKindName() === 'StringLiteral') {
        return (nameNode as import('ts-morph').StringLiteral).getLiteralValue()
      }
      return nameNode.getText()
    }),
  )

  const receivedKeySet = new Set(
    Object.keys(received).filter((key) => !ignoredKeys.includes(key)),
  )

  const unexpectedKeys = [...receivedKeySet].filter(
    (key) => !schemaKeySet.has(key),
  )
  const missingKeys = [...schemaKeySet].filter(
    (key) => !receivedKeySet.has(key),
  )

  const pass = unexpectedKeys.length === 0 && missingKeys.length === 0

  return {
    pass,
    message: () =>
      pass
        ? `Expected object not to match schema`
        : `Expected object to match schema
want:        ${[...schemaKeySet].join(', ')}
got:         ${[...receivedKeySet].join(', ')}
unexpected: ${unexpectedKeys.join(', ')}`,
  }
}

export interface SchemaVitestMatchers<R = unknown> {
  toMatchSchema(
    sourceFilePath: string,
    interfaceName: string,
    ignoredKeys?: string[],
  ): R
}

/**
 * Setup custom CDKTF Output Schema matchers for vitest.
 */
export async function setupSchemaMatchers() {
  // Dynamically import vitest
  const { expect } = await import('vitest')
  // Extends vitest's expect types.
  // See https://vitest.dev/guide/extending-matchers.html
  expect.extend({
    toMatchSchema,
  })
}
