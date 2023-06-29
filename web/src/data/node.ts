import useSWR from "swr";
import fetch from "../libs/fetch";
import { Node } from "../types/node";

export const useNode = ({
  nodeId,
}: {
  nodeId: string;
}): {
  loading: boolean;
  loggedOut: boolean;
  node: Node;
} => {
  const { data, error } = useSWR(
    `http://localhost:8080/v1/nodes/${nodeId}`,
    fetch
  );

  const loading = !data && !error;
  const loggedOut = error && error.status === 403;

  if (loggedOut) {
    window.location.href = "http://localhost:8080/v1/auth/";
  }

  return {
    loading,
    loggedOut,
    node: data as Node,
  };
};
