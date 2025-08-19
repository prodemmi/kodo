import { Item, ItemContext, OpenFileResponse } from "../types/item";

export const getItems = async (): Promise<Item[]> => {
  const response = await fetch(`http://localhost:8080/api/items`);
  return await response.json();
};

export const getItem = async (id: number): Promise<Item> => {
  const response = await fetch(`http://localhost:8080/api/items/${id}`);
  return await response.json();
};

export const getItemContext = async (
  file: string,
  line: number
): Promise<ItemContext> => {
  const response = await fetch(`http://localhost:8080/api/get-context`, {
    method: "POST",
    body: JSON.stringify({ file, line }),
  });
  return await response.json();
};

export const updateItem = async (
  id: number,
  status: string
): Promise<OpenFileResponse> => {
  const response = await fetch(`http://localhost:8080/api/items/update`, {
    method: "PUT",
    body: JSON.stringify({ id, status }),
  });
  return await response.json();
};

export const openFile = async (
  file: string,
  line: number
): Promise<OpenFileResponse> => {
  const response = await fetch(`http://localhost:8080/api/open-file`, {
    method: "POST",
    body: JSON.stringify({ file, line }),
  });
  return await response.json();
};
