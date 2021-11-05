<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/publicRoom/${page.params['room']}`);

		if (res.ok) {
			return {
				props: {
					room: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load that room`)
		};
	}
</script>

<script lang="ts">
	import type { PublicRoom } from '$lib/types/Room';

	export let room: PublicRoom;

	let selectedOption = null;
</script>

<div class="container mx-auto max-w-lg py-4">
	<form class="card shadow-lg">
		<div class="card-body">
			<div class="card-title">
				{room.question}
			</div>
			<div class="flex flex-col space-y-2">
				{#each room.options as option, i}
					<button
						class:btn-active={selectedOption === i}
						class:white={selectedOption === i}
						on:click={() => (selectedOption = i)}
						type="button"
						class="btn btn-secondary btn-outline border-2 not-uppercase"
					>
						{option}
					</button>
				{/each}
			</div>
			<div class="flex justify-end mt-4">
				<button class="btn btn-primary">Save</button>
			</div>
		</div>
	</form>
</div>

<style lang="postcss">
	.white {
		@apply text-white;
	}

	.not-uppercase {
		text-transform: none;
	}
</style>
