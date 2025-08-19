import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Item,
  ItemContext,
  OpenFileParams,
  OpenFileResponse,
  UpdateItemParams,
  UpdateItemResponse,
} from "../types/item";
import {
  getItem,
  getItemContext,
  getItems,
  openFile,
  updateItem,
} from "../api/item.api";

export function useItems() {
  return useQuery<Item[], Error>({
    queryKey: ["items"],
    queryFn: getItems,
  });
}

export function useItem(id: number) {
  return useQuery<Item, Error>({
    queryKey: ["item", id],
    queryFn: () => getItem(id),
    enabled: !!id,
  });
}

export function useItemContext(filename: string, line: number) {
  return useQuery<ItemContext, Error>({
    queryKey: ["item", "context", filename, line],
    queryFn: () => getItemContext(filename, line),
    enabled: !!filename && !!line,
  });
}

export function useOpenFile() {
  return useMutation<OpenFileResponse, Error, OpenFileParams>({
    mutationFn: ({ filename, line }) => openFile(filename, line),
  });
}

export function useUpdateItem() {
  const queryClient = useQueryClient();

  return useMutation<UpdateItemResponse, Error, UpdateItemParams>({
    mutationKey: ["items"],
    mutationFn: ({ id, status }) => updateItem(id, status),

    onSuccess: (_, variables) => {
      queryClient.setQueryData<Item[]>(["items"], (oldItems) =>
        oldItems?.map((item) =>
          item.id === variables.id
            ? { ...item }
            : item
        )
      );
    },
  });
}
