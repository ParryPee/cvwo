import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { fetchAllTopics } from "../api/forum"; //
import type { Topic } from "../types/models"; //

import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardActions from "@mui/material/CardActions";
import Button from "@mui/material/Button";

const ExploreTopicsPage = () => {
	const [topics, setTopics] = useState<Topic[]>([]);

	useEffect(() => {
		const loadTopics = async () => {
			try {
				const data = await fetchAllTopics();
				data.sort((a, b) => a.title.localeCompare(b.title));
				setTopics(data);
			} catch (error) {
				console.error("Failed to fetch topics", error);
			}
		};
		loadTopics();
	}, []);

	return (
		<Container maxWidth="lg" sx={{ mt: 4 }}>
			<Typography variant="h4" component="h1" gutterBottom>
				Explore Topics
			</Typography>
			<Grid container spacing={3}>
				{topics.map((topic) => (
					<Grid size={{ xs: 12, sm: 6, md: 4 }} key={topic.id}>
						<Card
							sx={{
								height: "100%",
								display: "flex",
								flexDirection: "column",
							}}
						>
							<CardContent sx={{ flexGrow: 1 }}>
								<Typography variant="h5" component="h2">
									{topic.title}
								</Typography>
								<Typography
									variant="body2"
									color="text.secondary"
									sx={{ mt: 1 }}
								>
									{topic.description}
								</Typography>
								<Typography
									variant="caption"
									display="block"
									sx={{ mt: 2 }}
								>
									Created by: {topic.created_by_username}
								</Typography>
							</CardContent>
							<CardActions>
								<Button
									size="small"
									component={Link}
									to={`/topics/${topic.id}`}
								>
									View Posts
								</Button>
							</CardActions>
						</Card>
					</Grid>
				))}
			</Grid>
		</Container>
	);
};

export default ExploreTopicsPage;
