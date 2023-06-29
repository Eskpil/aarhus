import React, { useEffect, useState } from "react";
import { useServer } from "../../data/server";
import { useParams } from "react-router-dom";
import {
  Box,
  Divider,
  Heading,
  HStack,
  Stat,
  StatLabel,
  StatNumber,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from "@chakra-ui/react";
import { DashboardLayout } from "../../layout/dashboard";
import { OverviewTab } from "./tabs/overview";
import { ConsoleTab } from "./tabs/console";
import { Server } from "../../types/server";

export const ServerPage: React.FC<{}> = () => {
  const params = useParams<{ id: string }>();

  const [selectedTab, setSelectedTab] = useState<number>(0);

  const updateSelectedTab = (server: Server, index: number): void => {
    localStorage.setItem(`server.${server.id}.tabs`, String(index));
    setSelectedTab(index);
  };

  const { server, loading } = useServer({ serverId: params.id! });

  useEffect(() => {
    if (server) {
      const storedTab = localStorage.getItem(`server.${server.id}.tabs`);
      if (storedTab) {
        setSelectedTab(parseInt(storedTab));
      }
    }
  }, [server]);

  if (loading) {
    return <div>loading......</div>;
  }

  return (
    <DashboardLayout>
      <Box bg="brand.700" padding={4} borderRadius="md">
        <HStack>
          <Heading>{server?.name}</Heading>
        </HStack>
      </Box>
      <Divider bg="brand.400" mt={4} />
      <Box mt={4}>
        <HStack>
          <Stat bg="red.300" padding={8} borderRadius="md">
            <StatLabel>Ram</StatLabel>
            <StatNumber>80% / {server.allocation.ram} GB</StatNumber>
          </Stat>
          <Stat bg="green.300" padding={8} borderRadius="md">
            <StatLabel>CPU</StatLabel>
            <StatNumber>55% / {server.allocation.cpu} vCPU</StatNumber>
          </Stat>
          <Stat bg="green.300" padding={8} borderRadius="md">
            <StatLabel>Players</StatLabel>
            <StatNumber>5 / 50</StatNumber>
          </Stat>
        </HStack>
      </Box>
      <Divider bg="brand.400" mt={4} />
      <Box mt={4}>
        <Tabs
          isFitted
          variant="enclosed"
          index={selectedTab}
          onChange={(index) => updateSelectedTab(server, index)}
        >
          <TabList>
            <Tab bg={selectedTab == 0 ? "brand.700" : "brand.100"}>
              <Heading as="h4" color="white" size="md">
                Overview
              </Heading>
            </Tab>
            <Tab bg={selectedTab == 1 ? "brand.700" : "brand.100"}>
              <Heading as="h4" color="white" size="md">
                Logs
              </Heading>
            </Tab>
            <Tab bg={selectedTab == 2 ? "brand.700" : "brand.100"}>
              <Heading as="h4" color="white" size="md">
                Console
              </Heading>
            </Tab>
          </TabList>
          <TabPanels>
            <TabPanel bg="brand.700" borderRadius="md">
              <OverviewTab server={server} />
            </TabPanel>
            <TabPanel bg="brand.700" borderRadius="md">
              <p>two!</p>
            </TabPanel>
            <TabPanel bg="brand.700" borderRadius="md">
              <ConsoleTab server={server} />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
    </DashboardLayout>
  );
};
