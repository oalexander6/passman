'use client';

import { useState } from 'react';
import {
    Dialog,
    DialogBackdrop,
    DialogPanel,
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    TransitionChild,
} from '@headlessui/react';
import {
    Cog6ToothIcon,
    FolderIcon,
    XMarkIcon,
    PencilSquareIcon
} from '@heroicons/react/24/outline';
import {
    ChevronRightIcon,
    ChevronUpDownIcon,
} from '@heroicons/react/20/solid';
import { classNames } from '../utils';

const navigation = [
    { name: 'Passwords', href: '#', icon: FolderIcon, current: true },
    { name: 'Notes', href: '#', icon: PencilSquareIcon, current: false },
    { name: 'Settings', href: '#', icon: Cog6ToothIcon, current: false },
];

const teams = [
    { id: 1, name: 'Planetaria', href: '#', initial: 'P', current: false },
    { id: 2, name: 'Protocol', href: '#', initial: 'P', current: false },
    { id: 3, name: 'Tailwind Labs', href: '#', initial: 'T', current: false },
];

const deployments = [
    {
        id: 1,
        href: '#',
        projectName: 'ios-app',
        teamName: 'Planetaria',
        status: 'offline',
        statusText: 'Initiated 1m 32s ago',
        description: 'Deploys from GitHub',
        environment: 'Preview',
    },
];

export default function BaseLayout() {
    const [sidebarOpen, setSidebarOpen] = useState(false);

    //   const { isPending, isError, data, error } = useQuery({
    //     queryKey: ['hellomsg'],
    //     queryFn: () => axios.get('/api').then(res => res.data)
    // });

    return (
        <>
            <div>
                <Dialog
                    open={sidebarOpen}
                    onClose={setSidebarOpen}
                    className="relative z-50 xl:hidden"
                >
                    <DialogBackdrop
                        transition
                        className="fixed inset-0 bg-gray-900/80 transition-opacity duration-300 ease-linear data-[closed]:opacity-0"
                    />

                    <div className="fixed inset-0 flex">
                        <DialogPanel
                            transition
                            className="relative mr-16 flex w-full max-w-xs flex-1 transform transition duration-300 ease-in-out data-[closed]:-translate-x-full"
                        >
                            <TransitionChild>
                                <div className="absolute left-full top-0 flex w-16 justify-center pt-5 duration-300 ease-in-out data-[closed]:opacity-0">
                                    <button
                                        type="button"
                                        onClick={() => setSidebarOpen(false)}
                                        className="-m-2.5 p-2.5"
                                    >
                                        <span className="sr-only">
                                            Close sidebar
                                        </span>
                                        <XMarkIcon
                                            aria-hidden="true"
                                            className="h-6 w-6 text-white"
                                        />
                                    </button>
                                </div>
                            </TransitionChild>
                            {/* Sidebar component, swap this element with another sidebar if you like */}
                            <div className="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 ring-1 ring-white/10">
                                <div className="flex h-16 shrink-0 items-center">
                                    <img
                                        alt="Your Company"
                                        src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=500"
                                        className="h-8 w-auto"
                                    />
                                </div>
                                <nav className="flex flex-1 flex-col">
                                    <ul
                                        role="list"
                                        className="flex flex-1 flex-col gap-y-7"
                                    >
                                        <li>
                                            <ul
                                                role="list"
                                                className="-mx-2 space-y-1"
                                            >
                                                {navigation.map((item) => (
                                                    <li key={item.name}>
                                                        <a
                                                            href={item.href}
                                                            className={classNames(
                                                                item.current
                                                                    ? 'bg-gray-800 text-white'
                                                                    : 'text-gray-400 hover:bg-gray-800 hover:text-white',
                                                                'group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6'
                                                            )}
                                                        >
                                                            <item.icon
                                                                aria-hidden="true"
                                                                className="h-6 w-6 shrink-0"
                                                            />
                                                            {item.name}
                                                        </a>
                                                    </li>
                                                ))}
                                            </ul>
                                        </li>
                                        <li>
                                            <div className="text-xs font-semibold leading-6 text-gray-400">
                                                Your teams
                                            </div>
                                            <ul
                                                role="list"
                                                className="-mx-2 mt-2 space-y-1"
                                            >
                                                {teams.map((team) => (
                                                    <li key={team.name}>
                                                        <a
                                                            href={team.href}
                                                            className={classNames(
                                                                team.current
                                                                    ? 'bg-gray-800 text-white'
                                                                    : 'text-gray-400 hover:bg-gray-800 hover:text-white',
                                                                'group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6'
                                                            )}
                                                        >
                                                            <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-[0.625rem] font-medium text-gray-400 group-hover:text-white">
                                                                {team.initial}
                                                            </span>
                                                            <span className="truncate">
                                                                {team.name}
                                                            </span>
                                                        </a>
                                                    </li>
                                                ))}
                                            </ul>
                                        </li>
                                        <li className="-mx-6 mt-auto">
                                            <a
                                                href="#"
                                                className="flex items-center gap-x-4 px-6 py-3 text-sm font-semibold leading-6 text-white hover:bg-gray-800"
                                            >
                                                <img
                                                    alt=""
                                                    src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                                                    className="h-8 w-8 rounded-full bg-gray-800"
                                                />
                                                <span className="sr-only">
                                                    Your profile
                                                </span>
                                                <span aria-hidden="true">
                                                    Tom Cook
                                                </span>
                                            </a>
                                        </li>
                                    </ul>
                                </nav>
                            </div>
                        </DialogPanel>
                    </div>
                </Dialog>

                {/* Static sidebar for desktop */}
                <div className="hidden xl:fixed xl:inset-y-0 xl:z-50 xl:flex xl:w-72 xl:flex-col">
                    {/* Sidebar component, swap this element with another sidebar if you like */}
                    <div className="flex grow flex-col gap-y-5 overflow-y-auto bg-black/10 px-6 ring-1 ring-white/5">
                        <div className="flex h-16 shrink-0 items-center">
                            <img
                                alt="Your Company"
                                src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=500"
                                className="h-8 w-auto"
                            />
                        </div>
                        <nav className="flex flex-1 flex-col">
                            <ul
                                role="list"
                                className="flex flex-1 flex-col gap-y-7"
                            >
                                <li>
                                    <ul role="list" className="-mx-2 space-y-1">
                                        {navigation.map((item) => (
                                            <li key={item.name}>
                                                <a
                                                    href={item.href}
                                                    className={classNames(
                                                        item.current
                                                            ? 'bg-gray-800 text-white'
                                                            : 'text-gray-400 hover:bg-gray-800 hover:text-white',
                                                        'group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6'
                                                    )}
                                                >
                                                    <item.icon
                                                        aria-hidden="true"
                                                        className="h-6 w-6 shrink-0"
                                                    />
                                                    {item.name}
                                                </a>
                                            </li>
                                        ))}
                                    </ul>
                                </li>
                                <li>
                                    <div className="text-xs font-semibold leading-6 text-gray-400">
                                        Your teams
                                    </div>
                                    <ul
                                        role="list"
                                        className="-mx-2 mt-2 space-y-1"
                                    >
                                        {teams.map((team) => (
                                            <li key={team.name}>
                                                <a
                                                    href={team.href}
                                                    className={classNames(
                                                        team.current
                                                            ? 'bg-gray-800 text-white'
                                                            : 'text-gray-400 hover:bg-gray-800 hover:text-white',
                                                        'group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6'
                                                    )}
                                                >
                                                    <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-[0.625rem] font-medium text-gray-400 group-hover:text-white">
                                                        {team.initial}
                                                    </span>
                                                    <span className="truncate">
                                                        {team.name}
                                                    </span>
                                                </a>
                                            </li>
                                        ))}
                                    </ul>
                                </li>
                                <li className="-mx-6 mt-auto">
                                    <a
                                        href="#"
                                        className="flex items-center gap-x-4 px-6 py-3 text-sm font-semibold leading-6 text-white hover:bg-gray-800"
                                    >
                                        <img
                                            alt=""
                                            src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                                            className="h-8 w-8 rounded-full bg-gray-800"
                                        />
                                        <span className="sr-only">
                                            Your profile
                                        </span>
                                        <span aria-hidden="true">Tom Cook</span>
                                    </a>
                                </li>
                            </ul>
                        </nav>
                    </div>
                </div>

                <div className="xl:pl-72">
                    <main className="lg:pr-96">
                        <header className="flex items-center justify-between border-b border-white/5 px-4 py-4 sm:px-6 sm:py-6 lg:px-8">
                            <h1 className="text-base font-semibold leading-7 text-white">
                                Deployments
                            </h1>

                            {/* Sort dropdown */}
                            <Menu as="div" className="relative">
                                <MenuButton className="flex items-center gap-x-1 text-sm font-medium leading-6 text-white">
                                    Sort by
                                    <ChevronUpDownIcon
                                        aria-hidden="true"
                                        className="h-5 w-5 text-gray-500"
                                    />
                                </MenuButton>
                                <MenuItems
                                    transition
                                    className="absolute right-0 z-10 mt-2.5 w-40 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 transition focus:outline-none data-[closed]:scale-95 data-[closed]:transform data-[closed]:opacity-0 data-[enter]:duration-100 data-[leave]:duration-75 data-[enter]:ease-out data-[leave]:ease-in"
                                >
                                    <MenuItem>
                                        <a
                                            href="#"
                                            className="block px-3 py-1 text-sm leading-6 text-gray-900 data-[focus]:bg-gray-50"
                                        >
                                            Name
                                        </a>
                                    </MenuItem>
                                    <MenuItem>
                                        <a
                                            href="#"
                                            className="block px-3 py-1 text-sm leading-6 text-gray-900 data-[focus]:bg-gray-50"
                                        >
                                            Date updated
                                        </a>
                                    </MenuItem>
                                    <MenuItem>
                                        <a
                                            href="#"
                                            className="block px-3 py-1 text-sm leading-6 text-gray-900 data-[focus]:bg-gray-50"
                                        >
                                            Environment
                                        </a>
                                    </MenuItem>
                                </MenuItems>
                            </Menu>
                        </header>

                        {/* Deployment list */}
                        <ul role="list" className="divide-y divide-white/5">
                            {deployments.map((deployment) => (
                                <li
                                    key={deployment.id}
                                    className="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8"
                                >
                                    <div className="min-w-0 flex-auto">
                                        <div className="flex items-center gap-x-3">
                                            <div className="text-green-400 bg-green-400/10 flex-none rounded-full p-1">
                                                <div className="h-2 w-2 rounded-full bg-current" />
                                            </div>
                                            <h2 className="min-w-0 text-sm font-semibold leading-6 text-white">
                                                <a
                                                    href={deployment.href}
                                                    className="flex gap-x-2"
                                                >
                                                    <span className="truncate">
                                                        {deployment.teamName}
                                                    </span>
                                                    <span className="text-gray-400">
                                                        /
                                                    </span>
                                                    <span className="whitespace-nowrap">
                                                        {deployment.projectName}
                                                    </span>
                                                    <span className="absolute inset-0" />
                                                </a>
                                            </h2>
                                        </div>
                                        <div className="mt-3 flex items-center gap-x-2.5 text-xs leading-5 text-gray-400">
                                            <p className="truncate">
                                                {deployment.description}
                                            </p>
                                            <svg
                                                viewBox="0 0 2 2"
                                                className="h-0.5 w-0.5 flex-none fill-gray-300"
                                            >
                                                <circle r={1} cx={1} cy={1} />
                                            </svg>
                                            <p className="whitespace-nowrap">
                                                {deployment.statusText}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="text-indigo-400 bg-indigo-400/10 ring-indigo-400/30 flex-none rounded-full px-2 py-1 text-xs font-medium ring-1 ring-inset">
                                        {deployment.environment}
                                    </div>
                                    <ChevronRightIcon
                                        aria-hidden="true"
                                        className="h-5 w-5 flex-none text-gray-400"
                                    />
                                </li>
                            ))}
                        </ul>
                    </main>
                </div>
            </div>
        </>
    );
}
