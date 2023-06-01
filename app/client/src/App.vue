<template>
  <ul>
    <li v-for="product of result.data.products" :key="product.id">
      {{ product.name }}
    </li>
  </ul>
  <table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Price</th>
        <th>Rank</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>Example Product</td>
        <td>5.00</td>
        <td>1</td>
      </tr>
    </tbody>
  </table>
</template>

<script setup lang="ts">
import { useQuery } from "@vue/apollo-composable";
import { gql } from "graphql-tag";

const getUsersQuery = gql`
  query {
    products(order: { field: "rank", direction: "ASC" }, page: { first: 10 }) {
      Total
      Edges {
        Cursor
        Node {
          id
          name
          price
          rank
        }
      }
    }
  }
`;

const { result } = useQuery(getUsersQuery);
</script>
<style scoped></style>
