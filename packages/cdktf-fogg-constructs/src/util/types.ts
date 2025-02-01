// helper type to make readonly interface properties mutable
export type Mutable<T> = {
  -readonly [P in keyof T]: T[P];
};

/**
 * Possible Terraform output types
 * that map to a known remote state getter.
 */
export type OutputType = "string" | "list" | "number" | "boolean";

/**
 * A small structure describing each output key
 * and which underlying `DataTerraformRemoteState` getter should be used
 */
export type OutputSchema<T> = {
  [K in keyof T]: OutputType;
};
