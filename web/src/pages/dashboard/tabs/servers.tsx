import React from "react";
import {
  Skeleton,
  Table,
  TableCaption,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import { useServers } from "../../../data/server";
import { useNode } from "../../../data/node";
import { Server } from "../../../types/server";

interface TableEntryProps {
  server: Server;
}

const TableEntry: React.FC<TableEntryProps> = ({ server }) => {
  const navigate = useNavigate();

  const { node, loading } = useNode({ nodeId: server.node_ref });
  if (loading) {
    return <Skeleton height="20px" />;
  }

  return (
    <Tr onClick={() => navigate(`/server/${server.id}`)}>
      <Td>{server.name}</Td>
      <Td>{node.name}</Td>
      <Td>{server.allocation.ram}</Td>
      <Td>{server.allocation.cpu}</Td>
      <Td>{server.last_started.toLocaleString()}</Td>
    </Tr>
  );
};

export const Servers: React.FC<{}> = () => {
  const { servers, loading } = useServers();
  if (loading) {
    return <div>loading...</div>;
  }

  return (
    <div>
      <TableContainer bg="brand.700" borderRadius="md">
        <Table variant="simple">
          <TableCaption>Servers</TableCaption>
          <Thead>
            <Tr>
              <Th>Name</Th>
              <Th>Node</Th>
              <Th>Ram</Th>
              <Th>Cpu</Th>
              <Th>Up since</Th>
            </Tr>
          </Thead>
          <Tbody>
            {servers.map((server) => (
              <TableEntry server={server} />
            ))}
          </Tbody>
        </Table>
      </TableContainer>
    </div>
  );
};
