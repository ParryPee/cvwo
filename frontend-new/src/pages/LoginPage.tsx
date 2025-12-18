import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
	Container,
	Box,
	Typography,
	TextField,
	Button,
	Alert,
	Paper,
} from "@mui/material";

const LoginPage: React.FC = () => {
	const { login } = useAuth();

	const navigate = useNavigate();

	const [username, setUsername] = useState("");
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault(); // Stop the page from reloading
		setError("");
		setLoading(true);

		try {
			await login(username);

			navigate("/");
		} catch (err) {
			console.error(err);
			setError("Failed to log in. Please try again.");
		} finally {
			setLoading(false);
		}
	};

	return (
		<Container maxWidth="xs">
			<Box
				sx={{
					marginTop: 8,
					display: "flex",
					flexDirection: "column",
					alignItems: "center",
				}}
			>
				<Paper elevation={3} sx={{ padding: 4, width: "100%" }}>
					<Typography
						component="h1"
						variant="h5"
						align="center"
						gutterBottom
					>
						Welcome Back
					</Typography>

					{error && (
						<Alert severity="error" sx={{ mb: 2 }}>
							{error}
						</Alert>
					)}

					<Box component="form" onSubmit={handleSubmit}>
						<TextField
							margin="normal"
							required
							fullWidth
							id="username"
							label="Username"
							name="username"
							autoFocus
							value={username}
							onChange={(e) => setUsername(e.target.value)}
							disabled={loading}
						/>

						<Button
							type="submit"
							fullWidth
							variant="contained"
							sx={{ mt: 3, mb: 2 }}
							disabled={loading}
						>
							{loading ? "Signing in..." : "Sign In"}
						</Button>
					</Box>
				</Paper>
			</Box>
		</Container>
	);
};

export default LoginPage;
