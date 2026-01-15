import { Grid, Fade, Box, Typography } from "@mui/material";

import DisplayCard from "./DisplayCard";

import type { Topic, Post } from "../types/models";

interface FeedProps {
	items: (Topic | Post)[];
	emptyMessage?: string;
}

const Feed = ({ items, emptyMessage = "No results found." }: FeedProps) => {
	if (!items || items.length === 0) {
		return (
			<Box sx={{ textAlign: "center", py: 8 }}>
				<Typography variant="h6" color="text.secondary">
					{emptyMessage}
				</Typography>
			</Box>
		);
	}
	return (
		<Grid container spacing={2}>
			{items.map((item, index) => {
				const isTopic = "description" in item;

				const title = item.title;

				const createdAt = item.created_at;
				const previewText = isTopic
					? (item as Topic).description
					: (item as Post).content;
				const linkTo = isTopic
					? `/topics/${item.id}`
					: `/topics/${(item as Post).topic_id}/posts/${item.id}`;
				return (
					<Fade
						key={item.id}
						in={true}
						timeout={800}
						style={{ transitionDelay: `${index * 100}ms` }}
					>
						<Grid size={{ xs: 12, sm: 6 }}>
							<DisplayCard
								title={title}
								previewText={previewText}
								createdAt={createdAt}
								linkTo={linkTo}
								topicTitle={!isTopic ? item.topic_title : ""}
								topicID={!isTopic ? item.topic_id : undefined}
							/>
						</Grid>
					</Fade>
				);
			})}
		</Grid>
	);
};

export default Feed;
