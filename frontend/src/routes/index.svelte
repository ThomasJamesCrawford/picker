<script lang="ts">
	import Alert from '$lib/icons/alert.svelte';
	import Check from '$lib/icons/check.svelte';
	import Plus from '$lib/icons/plus.svelte';
	import Trash from '$lib/icons/trash.svelte';
	import Info from '$lib/icons/info.svelte';
	import Loading from '$lib/icons/loading.svelte';
	import { debounce } from '$lib/helpers/debounce';
	import { nanoid } from 'nanoid';

	let options: string[] = [];

	let optionInputValue = '';
	let shortLink = nanoid(10);

	let shortLinkValidated = true;
	let shortLinkvalidationLoading = false;

	const addOption = () => {
		if (optionInputValue !== '') {
			options = [...options, optionInputValue];
			optionInputValue = '';
		}
	};

	const fetchIsAvailable = debounce(500, (name: string) =>
		fetch(`${import.meta.env.VITE_API_URL}/publicRoom/${name}/available`)
			.then((res) => res.json())
			.then((res) => {
				if (res.available === true) {
					shortLinkValidated = true;
				} else {
					shortLinkValidated = false;
				}
			})
			.catch((err) => {
				console.error(err);
			})
			.finally(() => {
				shortLinkvalidationLoading = false;
			})
	);

	$: if (shortLink) {
		shortLinkValidated = false;
		shortLinkvalidationLoading = true;
		fetchIsAvailable(shortLink);
	}
</script>

<div class="container mx-auto max-w-xl">
	<form class="p-4" on:submit|preventDefault="{}">
		<div class="form-control my-2">
			<label for="username" class="label justify-start space-x-2">
				<span class="label-text">Short link</span>
				<span
					data-tip="This is the url where your pickr will live!"
					class="text-gray-500 tooltip-right tooltip tooltip-info"
				>
					<Info />
				</span>
			</label>
			<label class="input-group">
				<span>pickr.com/</span>
				<input
					id="short_link"
					bind:value={shortLink}
					type="text"
					placeholder="honeylashes"
					class="input input-bordered w-full z-10"
				/>
				<button
					type="button"
					class:btn-success={shortLinkValidated}
					class:btn-error={!shortLinkValidated && !shortLinkvalidationLoading}
					class:btn-info={shortLinkvalidationLoading}
					class="btn no-animation"
				>
					{#if shortLinkValidated}
						<Check />
					{:else if shortLinkvalidationLoading}
						<Loading />
					{:else}
						<Alert />
					{/if}
				</button>
			</label>
			{#if !shortLinkValidated && !shortLinkvalidationLoading}
				<label for="short_link" class="label">
					<span class="label-text-alt">That short url is not available!</span>
				</label>
			{/if}
		</div>

		<div class="form-control my-2">
			<label for="username" class="label">
				<span class="label-text">Question</span>
			</label>
			<textarea
				class="textarea h-24 textarea-bordered w-full"
				placeholder="Which option would you like to select?"
			/>
		</div>

		<div class="form-control my-2">
			<label for="fef" class="label justify-start space-x-2">
				<span class="label-text">Add some options</span>
				<span
					data-tip="Each option can only ever be selected by one pickee (no duplicate selections)"
					class="text-gray-500 tooltip-right tooltip tooltip-info"
				>
					<Info />
				</span>
			</label>
			{#if options.length > 0}
				<div class="flex-col space-y-2 mb-2">
					{#each options as option, i}
						<div class="flex space-x-2">
							<div class="bg-gray-200 p-3 rounded-xl w-full">{option}</div>
							<button
								on:click={() => (options = [...options.slice(0, i), ...options.slice(i + 1)])}
								type="button"
								class="btn btn-accent btn-circle"><Trash /></button
							>
						</div>
					{/each}
				</div>
			{/if}
			<div class="flex space-x-2">
				<input
					on:keydown={(e) => {
						if (e.key === 'Enter') {
							addOption();
							e.preventDefault();
							e.stopPropagation();
						}
					}}
					bind:value={optionInputValue}
					type="text"
					placeholder="6:00 am Friday 25/06/21"
					class="w-full input input-bordered"
				/>
				<button on:click={addOption} type="button" class="btn btn-secondary btn-circle"
					><Plus /></button
				>
			</div>
		</div>

		<div class="mt-12 flex justify-end">
			<button
				disabled={options.length <= 0 || !shortLinkValidated}
				aria-disabled={options.length <= 0 || !shortLinkValidated}
				class="btn btn-primary">Create pickr</button
			>
		</div>
	</form>
</div>
