import { http } from "@/utils/http";

export type PageResult<T> = {
  list: T[];
  total: number;
  page: number;
  page_size: number;
};

export type City = {
  id: number;
  name: string;
  province: string;
  code: string;
  status: number;
};

export type Merchant = {
  id: number;
  name: string;
  city_id: number;
  contact_phone: string;
  address: string;
  status: number;
  rating: number;
};

export type Order = {
  id: number;
  order_no: string;
  user_id: number;
  merchant_id: number;
  pay_amount: number;
  status: number;
  created_at: string;
  items?: Array<{
    product_name: string;
    quantity: number;
    subtotal: number;
  }>;
};

export type Stats = {
  merchant_count: number;
  order_count: number;
  user_count: number;
};

export const getStats = () => http.get<Stats, never>("/api/v1/admin/stats");

export const getCities = () => http.get<City[], never>("/api/v1/admin/cities");

export const createCity = (data: object) =>
  http.post<City, object>("/api/v1/admin/cities", { data });

export const updateCityStatus = (id: number, status: number) =>
  http.request<City>("put", `/api/v1/admin/cities/${id}/status`, {
    data: { status }
  });

export const getMerchants = (params?: object) =>
  http.get<PageResult<Merchant>, object>("/api/v1/admin/merchants", {
    params
  });

export const updateMerchantStatus = (id: number, status: number) =>
  http.request<Merchant>("put", `/api/v1/admin/merchants/${id}/status`, {
    data: { status }
  });

export const getOrders = (params?: object) =>
  http.get<PageResult<Order>, object>("/api/v1/admin/orders", { params });
