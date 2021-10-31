<script lang="ts">
	import Alert from '$lib/icons/alert.svelte';
	import Check from '$lib/icons/check.svelte';
	import Plus from '$lib/icons/plus.svelte';
	import Trash from '$lib/icons/trash.svelte';
	import Info from '$lib/icons/info.svelte';
	import Loading from '$lib/icons/loading.svelte';
	import { debounce } from '$lib/helpers/debounce';
	import { customAlphabet } from 'nanoid';
	import { APP_NAME } from '$lib/constants';

	const nanoid = customAlphabet(
		'useandom26T198340PX75pxJACKVERYMINDBUSHWOLFGQZbfghjklqvwyzrict',
		10
	);

	let options: string[] = [];
	let question = '';

	let optionInputValue = '';
	let shortLink = nanoid();

	let shortLinkValidated = true;
	let shortLinkvalidationLoading = false;

	let submitLoading = false;

	let errorTooltipOpen = false;

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

	$: if (/^[a-zA-Z0-9]+$/.test(shortLink)) {
		errorTooltipOpen = false;
		shortLinkValidated = false;
		shortLinkvalidationLoading = true;
		fetchIsAvailable(shortLink);
	} else {
		shortLink = shortLink.replace(/[^a-zA-Z0-9]/gi, '');

		errorTooltipOpen = true;
		setTimeout(() => (errorTooltipOpen = false), 2000);
	}

	$: submit = () => {
		submitLoading = true;

		fetch(`${import.meta.env.VITE_API_URL}/room`, {
			method: 'POST',
			headers: {
				accepts: 'application/json'
			},
			body: JSON.stringify({ id: shortLink, options, question })
		}).finally(() => {
			submitLoading = false;
		});
	};
</script>

<div class="container mx-auto max-w-xl">
	<form class="p-4" on:submit|preventDefault={submit}>
		<div class="form-control my-2">
			<label for="username" class="label justify-start space-x-2">
				<span class="label-text">Short link</span>
				<span
					data-tip="This is the url where your {APP_NAME} will live"
					class="text-gray-500 tooltip-right tooltip tooltip-accent"
				>
					<Info />
				</span>
			</label>
			<span
				data-tip="Only letters and numbers are allowed"
				class="tooltip-accent tooltip-bottom"
				class:tooltip={errorTooltipOpen}
				class:tooltip-open={errorTooltipOpen}
			>
				<label class="input-group">
					<span>{APP_NAME}.com/</span>
					<input
						title="Only letters and numbers allowed"
						pattern="^[a-zA-Z0-9]*$"
						id="short_link"
						bind:value={shortLink}
						type="text"
						placeholder="myurl"
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
			</span>
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
				bind:value={question}
				required
				aria-required
				class="textarea h-24 textarea-bordered w-full"
				placeholder="Which option would you like to select?"
			/>
		</div>

		<div class="form-control my-2">
			<label for="fef" class="label justify-start space-x-2">
				<span class="label-text">Add some options</span>
				<span
					data-tip="Options to choose from"
					class="text-gray-500 tooltip-bottom tooltip tooltip-accent"
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
				class="btn btn-primary"
				class:loading={submitLoading}>Create {APP_NAME}</button
			>
		</div>
	</form>
</div>
