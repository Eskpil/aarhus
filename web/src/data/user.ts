import useSWR from "swr";
import { User } from "../types/user";
import fetch from "../libs/fetch";

export const useUser = (): {
  loading: boolean;
  loggedOut: boolean;
  user: User;
} => {
  const { data, error } = useSWR("http://localhost:8080/v1/@me/", fetch);

  const loading = !data && !error;
  const loggedOut = error && error.status === 403;

  if (loggedOut) {
    window.location.href = "http://localhost:8080/v1/auth/";
  }

  return {
    loading,
    loggedOut,
    user: data as User,
  };
};
