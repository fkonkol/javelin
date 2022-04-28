import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { createSlice } from "@reduxjs/toolkit";

const DOGS_API_KEY = "b2b33e8e-2444-4105-97f8-6508bff171fd";

interface Breed {
  id: string;
  name: string;
  image: {
    url: string
  };
};

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: 'https://api.thedogapi.com/v1',
    prepareHeaders(headers) {
      headers.set('x-api-key', DOGS_API_KEY);
      return headers;
    },
  }),
  endpoints(builder) {
    return {
      fetchBreed: builder.query<Breed[], number | void>({
        query(limit = 10) {
          return `/breeds?limit=${limit}` 
        },
      }),
    };
  },
});

export const { useFetchBreedQuery } = apiSlice;
