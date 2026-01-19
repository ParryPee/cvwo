import React from "react";
import ReactDOM from "react-dom/client";
import App from "./app";
import { AuthProvider } from "./context/AuthContext";
import "./index.css";
import { StrictMode } from "react";

const root = ReactDOM.createRoot(
	document.getElementById("root") as HTMLElement,
);

root.render(
	<React.StrictMode>
		<AuthProvider>
			<App />
		</AuthProvider>
	</React.StrictMode>,
);
