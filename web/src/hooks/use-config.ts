import { useQuery } from "@tanstack/react-query";
import { InvestorResponse } from "../types/config";
import { getInvestor } from "../api/config.api";

export function useInvestor() {
  return useQuery<InvestorResponse, Error>({
    queryKey: ["config", "investor"],
    queryFn: getInvestor,
  });
}
