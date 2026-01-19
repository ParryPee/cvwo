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
import { createPost } from "../api/forum";
import { useAuth } from "../context/AuthContext";

const CreatePostPage = () => {
	const navigate = useNavigate();
	const { topicId } = useParams<{ topicId: string }>();
	const { user } = useAuth();
	const [content, setContent] = useState("");
	const [title, setTitle] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		if (!topicId) return;
		if (!user) {
			setError("You must be logged in to create a post.");
			return;
		}
		const userId = user.id;
		try {
			setLoading(true);
			setError(null);

			await createPost({
				topic_id: parseInt(topicId, 10),
				title: title,
				content: content,
				user_id: userId,
			});

			navigate(`/topics/${topicId}`);
		} catch (err) {
			console.error(err);
			setError("Failed to create post. Please try again.");
		} finally {
			setLoading(false);
		}
	};

	return (
		<Container
			maxWidth={false}
			disableGutters
			sx={{
				backgroundImage:
					"linear-gradient(to bottom, var(--color-lavender-grey-700), var(--color-space-indigo-500))",
			}}
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
					Create a New Post
				</Typography>

				{error && (
					<Alert severity="error" sx={{ mb: 2 }}>
						{error}
					</Alert>
				)}

				<form onSubmit={handleSubmit}>
					<TextField
						label="Post Title"
						multiline
						rows={1}
						fullWidth
						variant="outlined"
						value={title}
						onChange={(e) => setTitle(e.target.value)}
						required
						sx={{ mb: 3 }}
						placeholder="Make sure your post is relevant to the topic!"
					/>
					<TextField
						label="Post Content"
						multiline
						rows={6}
						fullWidth
						variant="outlined"
						value={content}
						onChange={(e) => setContent(e.target.value)}
						required
						sx={{ mb: 3 }}
						placeholder="What's on your mind?"
					/>

					<Box display="flex" justifyContent="flex-end" gap={2}>
						<Button
							variant="outlined"
							onClick={() => navigate(-1)}
							disabled={loading}
							sx={{
								borderRadius: 4,
								bgcolor: "var(--color-flag-red-500)",
								":hover": {
									backgroundColor:
										"var(--color-flag-red-600)",
								},
								color: "white",
							}}
						>
							Cancel
						</Button>
						<Button
							type="submit"
							variant="contained"
							disabled={loading || !content.trim()}
							sx={{
								borderRadius: 4,
								bgcolor: "var(--color-flag-red-500)",
								":hover": {
									backgroundColor:
										"var(--color-flag-red-600)",
								},
								color: "white",
							}}
						>
							{loading ? "Posting..." : "Submit Post"}
						</Button>
					</Box>
				</form>
			</Paper>
		</Container>
	);
};

export default CreatePostPage;
