import { Node } from "./node";
import { ResourceAllocation } from "./resourceAllocation";

export interface Server {
  id: string;
  name: string;
  node_ref: string;
  node: Node;
  allocation: ResourceAllocation;
  last_started: Date;
}
