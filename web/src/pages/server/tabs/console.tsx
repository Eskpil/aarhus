import { Server } from "../../../types/server";
import { Box, Heading } from "@chakra-ui/react";

interface ConsoleTabProps {
  server: Server;
}

export const ConsoleTab: React.FC<ConsoleTabProps> = ({ server }) => {
  return (
    <Box>
      <Heading>Console!!!! {server.id}</Heading>
    </Box>
  );
};
