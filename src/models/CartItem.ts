import { Product } from "~/models/Product";

export type CartItem = {
  product: Product;
  count: number;
};

export type CartResponse = {
  cart: {
    items: CartItem[];
    status: string;
  };
  total: number;
};