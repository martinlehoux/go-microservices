import { ApiAuthenticationGateway } from "./gateways/authentication.gateway";
import { ApiHumanResourcesGateway } from "./gateways/human_resources.gateway";

export const humanResourcesGateway = new ApiHumanResourcesGateway();
export const authenticationGateway = new ApiAuthenticationGateway();
