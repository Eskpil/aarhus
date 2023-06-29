import React from "react";
import {
  Table,
  TableCaption,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";

export const Nodes: React.FC<{}> = () => {
  // TODO: Use zustand or something to make these breadcrumbs actually functional
  return (
    <div>
      <TableContainer>
        <Table variant="simple" colorScheme="purple">
          <TableCaption>Nodes</TableCaption>
          <Thead>
            <Tr>
              <Th>Name</Th>
              <Th>Location</Th>
              <Th>Last heartbeat</Th>
            </Tr>
          </Thead>
          <Tbody>
            <Tr>
              <Td>node1</Td>
              <Td>Frankfurt</Td>
              <Td>{new Date().toLocaleString()}</Td>
            </Tr>
            <Tr>
              <Td>node2</Td>
              <Td>Munich</Td>
              <Td>{new Date().toLocaleString()}</Td>
            </Tr>
            <Tr>
              <Td>node3</Td>
              <Td>Helsinki</Td>
              <Td>{new Date().toLocaleString()}</Td>
            </Tr>
          </Tbody>
        </Table>
      </TableContainer>
    </div>
  );
};
