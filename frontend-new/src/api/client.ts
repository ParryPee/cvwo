import axios from "axios";

const baseURL = import.meta.env.VITE_API_URL || "http://localhost:8080";

const instance = axios.create({
	baseURL: `${baseURL}/api/`,
	withCredentials: true,
});

export default instance;