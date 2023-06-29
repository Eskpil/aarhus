import useSWR from "swr";
import fetch from "../libs/fetch";
import { Server } from "../types/server";

export const useServer = ({
  serverId,
}: {
  serverId: string;
}): { loading: boolean; loggedOut: boolean; server: Server } => {
  const { data, error } = useSWR(
    `http://localhost:8080/v1/servers/${serverId}`,
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
    server: data as Server,
  };
};

export const useServers = (): {
  loading: boolean;
  loggedOut: boolean;
  servers: Server[];
} => {
  const { data, error } = useSWR("http://localhost:8080/v1/servers/", fetch);

  const loading = !data && !error;
  const loggedOut = error && error.status === 403;

  if (loggedOut) {
    window.location.href = "http://localhost:8080/v1/auth/";
  }

  return {
    loading,
    loggedOut,
    servers: data as Server[],
  };
};
