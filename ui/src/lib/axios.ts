import axios, { AxiosRequestConfig, AxiosResponse, AxiosError } from "axios";

const client = axios.create({
  baseURL: "http://127.0.0.1:8000",
  withCredentials: true,
  withXSRFToken: true,
});

export const request = (options: AxiosRequestConfig<any>) => {
  const onSuccess = (response: AxiosResponse<any, any>) => {
    console.log("axios response", response);
    console.log("axios response headers", response.headers.toString());
    return response;
  };
  const onError = (error: AxiosError<any, any>) => {
    console.log("axios error", error);
    if (error.response) {
      console.log("axios data", error.response.data);
      throw new Error(error.response.data.error);
    }

    throw error;
  };
  return client(options).then(onSuccess).catch(onError);
};
