<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/room`);

		if (res.ok) {
			return {
				props: {
					rooms: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load rooms`)
		};
	}
</script>

<script lang="ts">
	import type { SimpleRoom } from '$lib/types/Room';

	export let rooms: SimpleRoom[];
</script>

<div class="container mx-auto max-w-lg py-4">
	<div class="card bg-white shadow-lg">
		<div class="card-body">
			<div class="card-title">Rooms</div>
			<div>
				<table class="table w-full table-fixed">
					<thead>
						<tr>
							<th>name</th>
							<th>question</th>
						</tr>
					</thead>
					<tbody>
						{#each rooms as room}
							<tr>
								<td><a class="link link-secondary" href={`/admin/${room.id}`}>{room.id}</a></td>
								<td class="overflow-ellipsis overflow-hidden">{room.question}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>
