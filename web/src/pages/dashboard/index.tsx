import React from "react";
import { Divider, HStack, Stat, StatLabel, StatNumber } from "@chakra-ui/react";
import { Servers } from "./tabs/servers";
import { DashboardLayout } from "../../layout/dashboard";

export const DashboardPage: React.FC<{}> = () => {
  return (
    <DashboardLayout>
      <HStack>
        <Stat bg="green.300" padding={8} borderRadius="md">
          <StatLabel>Nodes</StatLabel>
          <StatNumber>3</StatNumber>
        </Stat>
        <Stat bg="orange.300" padding={8} borderRadius="md">
          <StatLabel>Servers</StatLabel>
          <StatNumber>12</StatNumber>
        </Stat>
        <Stat bg="red.300" padding={8} borderRadius="md">
          <StatLabel>Users</StatLabel>
          <StatNumber>24532</StatNumber>
        </Stat>
      </HStack>

      <Divider mt={8} />

      <Servers />
    </DashboardLayout>
  );
};
