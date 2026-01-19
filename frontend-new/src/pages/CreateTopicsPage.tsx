import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
	Container,
	Typography,
	TextField,
	Button,
	Box,
	Paper,
	Alert,
} from "@mui/material";
import { createTopic } from "../api/forum";
import { useAuth } from "../context/AuthContext";

const CreateTopicsPage = () => {
	const navigate = useNavigate();
	const { user } = useAuth();
	const [desc, setDesc] = useState("");
	const [title, setTitle] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		if (!user) {
			setError("You must be logged in to create a topic.");
			return;
		}
		const userId = user.id;
		try {
			setLoading(true);
			setError(null);
			await createTopic({
				title: title,
				description: desc,
				created_by: userId,
			});
			navigate(`/`);
		} catch (err) {
			console.error(err);
			setError("Failed to create topic. Please try again.");
		} finally {
			setLoading(false);
		}
	};

	return (
		<Container
			sx={{
				backgroundImage:
					"linear-gradient(to bottom, var(--color-lavender-grey-700), var(--color-space-indigo-500))",
			}}
			maxWidth={false}
			disableGutters
		>
			<Paper
				sx={{
					p: 4,
					maxWidth: "960px",
					margin: "0 auto",
					bgcolor: "var(--color-platinum-100)",
					borderRadius: 4,
					pt: 6,
				}}
			>
				<Typography variant="h5" gutterBottom>
					Create a New Topic
				</Typography>
				{error && (
					<Alert severity="error" sx={{ mb: 2 }}>
						{error}
					</Alert>
				)}
				<Box
					component="form"
					onSubmit={handleSubmit}
					sx={{ display: "flex", flexDirection: "column", gap: 2 }}
				>
					<TextField
						label="Title"
						value={title}
						onChange={(e) => setTitle(e.target.value)}
						required
					/>
					<TextField
						label="Description"
						value={desc}
						onChange={(e) => setDesc(e.target.value)}
						required
						multiline
						rows={4}
					/>
					<Button
						type="submit"
						variant="contained"
						sx={{
							backgroundColor: "var(--color-flag-red-500)",
							borderRadius: 4,
							":hover": {
								backgroundColor: "var(--color-flag-red-600)",
							},
							maxWidth: "150px",
							margin: "0 auto",
						}}
						disabled={loading}
					>
						{loading ? "Creating..." : "Create Topic"}
					</Button>
				</Box>
			</Paper>
		</Container>
	);
};
export default CreateTopicsPage;
