import React from "react";
import { Box } from "@chakra-ui/react";
import { Sidebar } from "./sidenav";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({
  children,
}) => {
  // TODO: Create menu

  return (
    <Sidebar>
      <Box padding={12} bg="brand.100" width="full">
        <Box>{children}</Box>
      </Box>
    </Sidebar>
  );
};
