import axios, {type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';

const request: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// 请求拦截器
request.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// 响应拦截器
request.interceptors.response.use(
    (response: AxiosResponse) => {
        // Axios 的 response.data 才是后端返回的内容
        // 这里我们可以根据后端返回的结构（比如 {code: 200, data: {}}）进行二次解构
        return response.data;
    },
    (error) => {
        // 处理 HTTP 状态码错误
        const message = error.response?.data?.error || error.message || 'Unknown Error';
        console.error('API Error:', message);
        return Promise.reject(error);
    }
);

export default request;
