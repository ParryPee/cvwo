import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
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
		<Container maxWidth="md" sx={{ mt: 4 }}>
			<Paper sx={{ p: 4 }}>
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
						color="primary"
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
