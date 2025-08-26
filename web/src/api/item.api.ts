import { Item, ItemContext, OpenFileResponse } from "../types/item";
import api from "../utils/api";

export const getItems = async (): Promise<Item[]> => {
  const response = await api.get<Item[]>("/items");
  return response.data;
};

export const getItem = async (id: number): Promise<Item> => {
  const response = await api.get<Item>(`/items/${id}`);
  return response.data;
};

export const getItemContext = async (
  file: string,
  line: number
): Promise<ItemContext> => {
  const response = await api.post<ItemContext>("/items/get-context", { file, line });
  return response.data;
};

export const updateItem = async (
  id: number,
  status: string
): Promise<OpenFileResponse> => {
  const response = await api.put<OpenFileResponse>("/items/update", {
    id,
    status,
  });
  return response.data;
};

export const openFile = async (
  file: string,
  line: number
): Promise<OpenFileResponse> => {
  const response = await api.post<OpenFileResponse>("/items/open-file", {
    file,
    line,
  });
  return response.data;
};
