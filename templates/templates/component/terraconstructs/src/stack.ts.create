import { FoggTerraStack } from "@vincenthsh/cdktf-fogg-helpers";
import { Construct } from "constructs";

export class ComponentStack extends FoggTerraStack {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      forceRemoteStates: false,
      // Prefix for resource UUID, should never change for component lifecycle
      gridUUID: "123-456",
      // This is used as a tag and may change over time
      environmentName: "development",
    });

    // For example usage see:
    // https://github.com/vincenthsh/fogg/blob/cdktf-fogg-helpers-v1.0.0/packages/cdktf-fogg-constructs
  }
}
