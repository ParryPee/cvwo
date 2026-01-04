import { useEffect, useState } from "react";
import { useNavigate, Link as RouterLink } from "react-router-dom";
import { fetchAllTopics } from "../api/forum";
import type { Topic } from "../types/models";
import { timeAgo, formatDate } from "../utils/date";
import Box from "@mui/material/Box";
import {
	Container,
	CircularProgress,
	Typography,
	Grid,
	Card,
	Button,
	CardActionArea,
	Fade,
} from "@mui/material";
import { Link, Padding } from "@mui/icons-material";
import { useAuth } from "../context/AuthContext";

const HomePage = () => {
	const navigate = useNavigate();
	const [topics, setTopics] = useState<Topic[]>([]);
	const [loading, setLoading] = useState(true);
	const { isAuthenticated, user } = useAuth();

	useEffect(() => {
		const loadTopics = async () => {
			try {
				const data = await fetchAllTopics();
				setTopics(data);
			} catch (error) {
				console.error("Failed to fetch topics:", error);
			} finally {
				setLoading(false);
			}
		};

		loadTopics();
	}, []);

	if (loading) {
		return (
			<Container>
				<CircularProgress />
			</Container>
		);
	} else {
		return (
			<Box
				sx={{
					border: "1px solid #ccc",
					borderRadius: "8px",
					padding: "16px",
					paddingX: "4rem",
					marginTop: "16px",
				}}
			>
				<Box
					sx={{
						display: "flex",
						alignItems: "center",
						justifyContent: "space-between",
						mb: 2,
					}}
				>
					<Typography variant="h4" gutterBottom sx={{ mb: 0 }}>
						Home Page
					</Typography>
					{isAuthenticated && (
						<Button
							variant="contained"
							color="primary"
							onClick={() => navigate(`/topics/create`)}
						>
							Create Topic
						</Button>
					)}
				</Box>
				{!topics || topics.length === 0 ? (
					<Box
						sx={{
							textAlign: "center",
							padding: "48px 16px",
						}}
					>
						<Typography variant="h6" color="text.secondary">
							No topics available yet
						</Typography>
						<Typography
							variant="body2"
							color="text.secondary"
							sx={{ mt: 1 }}
						>
							Be the first to create a topic!
						</Typography>
					</Box>
				) : (
					<Grid container spacing={2} sx={{ padding: "16px" }}>
						{topics.map((topic, index) => (
							<Fade
								key={topic.id}
								in={true}
								timeout={800}
								style={{ transitionDelay: `${index * 150}ms` }}
							>
								<Grid size={{ xs: 12, sm: 6 }}>
									<CardActionArea
										sx={{
											padding: "16px",
											border: "1px solid #ddd",
											borderRadius: "8px",
											display: "block",
										}}
										component={RouterLink}
										to={`/topics/${topic.id}`}
									>
										<Typography
											variant="h6"
											component="h3"
											gutterBottom
											sx={{
												fontWeight: 600,
											}}
										>
											{topic.title}
										</Typography>
										<Typography
											variant="body2"
											color="text.secondary"
											gutterBottom
											fontWeight={600}
										>
											Created {timeAgo(topic.created_at)}
										</Typography>
										<Typography
											variant="body2"
											color="text.secondary"
											sx={{
												display: "-webkit-box",
												WebkitLineClamp: 2,
												WebkitBoxOrient: "vertical",
												overflow: "hidden",
												textOverflow: "ellipsis",
											}}
										>
											{topic.description}
										</Typography>
									</CardActionArea>
								</Grid>
							</Fade>
						))}
					</Grid>
				)}
			</Box>
		);
	}
};

export default HomePage;
