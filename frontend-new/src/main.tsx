import React from "react";
import ReactDOM from "react-dom/client";
import App from "./app";
import { AuthProvider } from "./context/AuthContext";

const root = ReactDOM.createRoot(
	document.getElementById("root") as HTMLElement
);

root.render(
	<React.StrictMode>
		{/* The Electricity (Auth) wraps the entire house */}
		<AuthProvider>
			<App />
		</AuthProvider>
	</React.StrictMode>
);
