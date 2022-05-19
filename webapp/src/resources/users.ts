import { createResource } from "solid-js";
import { humanResourcesGateway } from "../module";

export type User = {
  id: string;
  preferred_name: string;
  email: string;
};

export const [users, { mutate, refetch }] = createResource(() =>
  humanResourcesGateway.listUsers()
);
