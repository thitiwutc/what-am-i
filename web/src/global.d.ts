export {};

export interface Env {
  env: string;
  apiUrl: string;
}

declare global {
  interface Window {
    __ENV: Env | undefined;
  }
}
