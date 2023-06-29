import React from "react";
import { Server } from "../../../types/server";

interface OverviewTabProps {
  server: Server;
}

export const OverviewTab: React.FC<OverviewTabProps> = ({ server }) => {
  return (
    <div>
      epic overview of {server.name}
      <div>some kind of terminal</div>
    </div>
  );
};
