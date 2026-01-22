import { useEffect, useState } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { fetchTopicById, fetchPostsByTopicId, updateTopic } from "../api/forum";
import type { Topic, Post } from "../types/models";
import { useAuth } from "../context/AuthContext";
import { timeAgo, formatDate } from "../utils/date";
import Box from "@mui/material/Box";
import EditModal from "../components/EditModal";
import {
	Container,
	CircularProgress,
	Typography,
	Grid,
	CardActionArea,
	Alert,
	Button,
	Divider,
	IconButton,
} from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";
const TopicPage = () => {
	const navigate = useNavigate();
	const { topicId } = useParams<{ topicId: string }>();
	const [topic, setTopic] = useState<Topic | null>(null);
	const [loading, setLoading] = useState(true);
	const [posts, setPosts] = useState<Post[]>([]);
	const [error, setError] = useState<string | null>(null);
	const [isEditModalOpen, setIsEditModalOpen] = useState(false);
	const { isAuthenticated, user } = useAuth();

	const handleUpdateTopic = async (title: string, description: string) => {
		if (!topic) return;
		try {
			await updateTopic(topic.id, { title, description });
			setTopic({ ...topic, title, description });
			setIsEditModalOpen(false);
		} catch (error) {
			console.error("Failed to update topic:", error);
		}
	};

	console.log(topic);
	useEffect(() => {
		const loadTopicAndPosts = async () => {
			if (!topicId) return;
			try {
				setLoading(true);
				const topicIdNum = parseInt(topicId, 10);
				if (isNaN(topicIdNum)) {
					setError("Invalid Topic ID");
					setLoading(false);
					return;
				}
				const topicData = await fetchTopicById(topicIdNum);
				const postsData = await fetchPostsByTopicId(topicIdNum);
				setPosts(postsData);
				setTopic(topicData);
			} catch (error) {
				console.error("Failed to fetch topic or posts:", error);
				setError("Failed to load topic. It may not exist.");
			} finally {
				setLoading(false);
			}
		};

		loadTopicAndPosts();
	}, [topicId]);
	if (loading) {
		return (
			<Box display="flex" justifyContent="center" mt={4}>
				<CircularProgress />
			</Box>
		);
	}
	if (error || !topic) {
		return (
			<Container sx={{ mt: 4 }}>
				<Alert severity="error">{error || "Topic not found"}</Alert>
			</Container>
		);
	}
	return (
		<Container
			maxWidth={false}
			disableGutters
			sx={{
				width: "100%",
				padding: "16px",
				backgroundImage:
					"linear-gradient(to bottom, var(--color-lavender-grey-700), var(--color-space-indigo-500))",
			}}
		>
			<Box
				sx={{
					border: "1px solid #ccc",
					borderRadius: 4,
					p: 2,
					bgcolor: "var(--color-platinum-300)",
					maxWidth: "960px",
					margin: "0 auto",
					mb: 4,
				}}
			>
				<Typography variant="h4" gutterBottom>
					{topic.title}
				</Typography>
				<Typography
					variant="subtitle1"
					color="text.secondary"
					gutterBottom
				>
					{topic.description}
				</Typography>
				<Box
					sx={{
						display: "flex",
						justifyContent: "flex-end",
					}}
				>
					{isAuthenticated && user?.id === topic.user_id && (
						<IconButton onClick={() => setIsEditModalOpen(true)}>
							<EditIcon
								sx={{
									cursor: "pointer",
									":hover": { color: "blue" },
								}}
							/>
						</IconButton>
					)}
				</Box>
			</Box>

			<Box mt={4} sx={{ maxWidth: "960px", margin: "0 auto", mt: 4 }}>
				<Box
					sx={{
						display: "flex",
						justifyContent: "space-between",
						alignItems: "center",
						mb: 2,
					}}
				>
					<Typography variant="h6">Posts</Typography>
					{isAuthenticated && (
						<Button
							variant="contained"
							sx={{
								backgroundColor: "var(--color-flag-red-500)",
								borderRadius: 4,
								":hover": {
									backgroundColor:
										"var(--color-flag-red-600)",
								},
							}}
							onClick={() =>
								navigate(`/topics/${topicId}/create`)
							}
						>
							Create Post
						</Button>
					)}
				</Box>
				<Divider sx={{ mb: 2 }} />
				{!posts || posts.length === 0 ? (
					<Typography variant="body2" color="text.secondary" mt={2}>
						No posts yet. Be the first to post!
					</Typography>
				) : (
					posts.map((post) => (
						<Grid
							key={post.id}
							size={{ xs: 12, sm: 6 }}
							sx={{
								mb: 2,
								bgcolor: "var(--color-platinum-200)",
								borderRadius: 4,
							}}
						>
							<CardActionArea
								sx={{
									border: "1px solid #ddd",
									borderRadius: 4,
									padding: "16px",
									marginBottom: "16px",
								}}
								component={Link}
								to={`/topics/${topicId}/posts/${post.id}`}
							>
								<Typography variant="h6">
									{post.title}
								</Typography>
								<Typography
									variant="body1"
									sx={{
										display: "-webkit-box",
										WebkitLineClamp: 2,
										WebkitBoxOrient: "vertical",
										overflow: "hidden",
										textOverflow: "ellipsis",
									}}
								>
									{post.content}
								</Typography>
								<Typography
									variant="caption"
									color="text.secondary"
								>
									Posted by • {post.created_by_username} •{" "}
									{timeAgo(post.created_at)} (on{" "}
									{formatDate(post.created_at)})
								</Typography>
								<Box my={2} />
							</CardActionArea>
						</Grid>
					))
				)}
			</Box>
			{topic && (
				<EditModal
					open={isEditModalOpen}
					onClose={() => setIsEditModalOpen(false)}
					onSubmit={handleUpdateTopic}
					initialTitle={topic.title}
					initialContent={topic.description}
				/>
			)}
		</Container>
	);
};
export default TopicPage;
