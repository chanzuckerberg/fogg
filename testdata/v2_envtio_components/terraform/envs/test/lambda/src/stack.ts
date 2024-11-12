// // Use provider-aws constructs directly
// import { dataAwsAvailabilityZones } from "@cdktf/provider-aws";
// // Or use CDKTF constructs directly from npm
// import { MyConstruct } from "@handshakes/my-construct";
import { Construct } from "constructs";
import { FoggStack } from "./helpers/fogg-stack";

export class ComponentStack extends FoggStack {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      forceRemoteStates: false,
      // Prefix for resource UUID, should never change for component lifecycle
      gridUUID: "123-456",
      // This is used as a tag and may change over time
      environmentName: "development",
    });

    // // Configure fogg module variables here
    // this.setModuleVariables("main", {
    //   foo: "bar",
    //   baz: "qux",
    //   tags: {
    //     // Use fogg component configuration values
    //     Project: this.foggComp.project,
    //     Owner: this.foggComp.owner,
    //     Environment: this.foggComp.env,
    //   },
    // });

    // // Or: create AWS Resources
    // const azs = new dataAwsAvailabilityZones.DataAwsAvailabilityZones(
    //   this,
    //   "azs",
    //   {}
    // );

    // // Or: use custom constructs
    // new MyConstruct(this, "my-construct", {
    //   foo: "bar",
    // });
  }
}
