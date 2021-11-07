import type { Option, PublicOption } from '$lib/types/Room';

export const sortOptions = <T extends PublicOption | Option>(options: T[]): T[] =>
	options.sort((a, b) => {
		if (a.value === b.value) {
			return a.id.localeCompare(b.id);
		}

		return a.value.localeCompare(b.value);
	});
