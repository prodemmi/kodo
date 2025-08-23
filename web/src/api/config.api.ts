import { InvestorResponse } from "../types/config";
import api from "../utils/api";

export const getInvestor = async (): Promise<InvestorResponse> => {
  try {
    const response = await api.get<InvestorResponse>("/investor");

    return response.data;
  } catch (error: any) {
    // Axios interceptor already shows notification
    throw new Error(
      error.response?.data?.message || "Failed to load project files"
    );
  }
};
