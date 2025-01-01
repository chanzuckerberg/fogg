// // Or publish cdktf module as HCL module for non CDKTF usage
// import { TFModuleStack, TFModuleVariable, ProviderRequirement } from "@cdktf/tf-module-stack";
import {
  TerraformElement,
  // TerraformHclModule,
} from "cdktf";
import { Construct } from "constructs";

export interface MyConstructProps {
  foo?: string;
}

/**
 * This is an example construct you may change as required.
 */
export class MyConstruct extends TerraformElement {
  private readonly foo: string;
  constructor(scope: Construct, id: string, props?: MyConstructProps) {
    super(scope, id);

    this.foo = props?.foo ?? "Default";

    // // create AWS Resources
    // const azs = new dataAwsAvailabilityZones.DataAwsAvailabilityZones(
    //   this,
    //   "azs",
    //   {}
    // );

    // // Or: use custom constructs
    // new MyConstruct(this, "my-construct", {
    //   foo: "bar",
    // });

    // // Or use community TF modules
    // https://developer.hashicorp.com/terraform/cdktf/create-and-deploy/configuration-file#declare-providers-and-modules
  }
}
