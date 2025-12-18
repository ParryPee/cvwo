
const parseDate = (dateString: string): Date => {
	// Replace space with T
	const isoString = dateString.trim().replace(" ", "T");
	const date = new Date(isoString);
	
	// Fallback
	return Number.isNaN(date.getTime()) ? new Date() : date;
};

export const formatDate = (isoString: string): string => {
	const date = parseDate(isoString);

	return new Intl.DateTimeFormat("en-US", {
		month: "short",
		day: "numeric",
		year: "numeric",
		hour: "numeric",
		minute: "numeric",
		hour12: true,
	}).format(date);
};

export const timeAgo = (isoString: string): string => {
	const date = parseDate(isoString);
	const now = new Date();

	const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

	let interval = Math.floor(seconds / 31536000);
	if (interval >= 1) return `${interval} ${interval === 1 ? "year" : "years"} ago`;

	interval = Math.floor(seconds / 2592000);
	if (interval >= 1) return `${interval} ${interval === 1 ? "month" : "months"} ago`;

	interval = Math.floor(seconds / 86400);
	if (interval >= 7) return `${Math.floor(interval / 7)} ${Math.floor(interval / 7) === 1 ? "week" : "weeks"} ago`;
	if (interval >= 1) return `${interval} ${interval === 1 ? "day" : "days"} ago`;

	interval = Math.floor(seconds / 3600);
	if (interval >= 1) return `${interval} ${interval === 1 ? "hour" : "hours"} ago`;

	interval = Math.floor(seconds / 60);
	if (interval >= 1) return `${interval} ${interval === 1 ? "minute" : "minutes"} ago`;

	return "just now";
};