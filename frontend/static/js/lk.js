const API = {
    profile: "/api/user/profile",
    projects: "/api/profile/projects",
    favorites: "/api/profile/favorites",
    updateProfile: "/api/user/profile", 
    changePassword: "/api/auth/password",
    logout: "/api/auth/logout"
};
    /*
     * ==========================================
     * DOM
     * ==========================================
     */

    const navItems = document.querySelectorAll(".profile-nav-item");

    const panels = document.querySelectorAll(".profile-panel");

    const projectsList =
        document.getElementById("projectsList");

    const projectsLoading =
        document.getElementById("projectsLoading");

    const projectsEmpty =
        document.getElementById("projectsEmpty");

    const favoritesList =
        document.getElementById("favoritesList");

    const favoritesLoading =
        document.getElementById("favoritesLoading");

    const favoritesEmpty =
        document.getElementById("favoritesEmpty");

    const profileForm =
        document.getElementById("profileForm");

    const passwordForm =
        document.getElementById("passwordForm");

    const logoutButton =
        document.getElementById("logoutButton");


    /*
     * ==========================================
     * API REQUEST
     * ==========================================
     */

async function apiRequest(url, options = {}) {
    const token = localStorage.getItem("access_token");

    const response = await fetch(url, {
        ...options,

        credentials: "include",

        headers: {
            "Content-Type": "application/json",

            ...(token
                ? {
                    "Authorization": `Bearer ${token}`
                }
                : {}),

            ...(options.headers || {})
        }
    });

    if (response.status === 401) {
        localStorage.removeItem("access_token");
        localStorage.removeItem("token");

        window.location.href = "/login";

        throw new Error("Unauthorized");
    }

    let data = null;

    const contentType = response.headers.get("content-type");

    if (
        contentType &&
        contentType.includes("application/json")
    ) {
        data = await response.json();
    }

    if (!response.ok) {
        const message =
            data?.message ||
            data?.error ||
            "Произошла ошибка";

        throw new Error(message);
    }

    return data;
}


    /*
     * ==========================================
     * PROFILE
     * ==========================================
     */

async function loadProfile() {
    try {
        const response = await fetch(API.profile, {
            credentials: "include",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("access_token")}`
            }
        });

        if (response.status === 401) {
            localStorage.removeItem("access_token");
            window.location.href = "/auth";
            return;
        }

        // Если профиль еще не создан (404), просто считаем его пустым, а не кидаем ошибку
        if (response.status === 404) {
            console.log("Профиль еще не создан, показываем пустые поля");
            const profile = { name: "", phone: "" };
            renderProfile(profile);
            return;
        }

        const profile = await response.json();
        renderProfile(profile);

    } catch (error) {
        console.error("Ошибка загрузки профиля:", error);
    }
}

// Вынеси отрисовку в отдельную функцию для чистоты
function renderProfile(profile) {
    const fullNameElement = document.getElementById("userFullName");
    if (fullNameElement) {
        fullNameElement.textContent = profile.name || "Пользователь";
    }

    const nameInput = document.getElementById("name");
    if (nameInput) {
        nameInput.value = profile.name || "";
    }
}

    /*
     * ==========================================
     * PROJECTS
     * ==========================================
     */

    async function loadProjects() {

        projectsLoading.classList.remove("hidden");

        projectsEmpty.classList.add("hidden");

        projectsList.innerHTML = "";


        try {

            const projects =
                await apiRequest(API.projects);


            projectsLoading.classList.add("hidden");


            if (
                !projects ||
                !Array.isArray(projects) ||
                projects.length === 0
            ) {

                projectsEmpty.classList.remove("hidden");

                return;
            }


            projects.forEach(project => {

                projectsList.appendChild(
                    createProjectCard(project)
                );

            });

        } catch (error) {

            projectsLoading.classList.add("hidden");

            console.error(
                "Ошибка загрузки проектов:",
                error
            );

            projectsEmpty.classList.remove("hidden");
        }
    }


    /*
     * ==========================================
     * PROJECT CARD
     * ==========================================
     */

    function createProjectCard(project) {

        const article =
            document.createElement("article");

        article.className =
            "project-card";


        /*
         * Все значения ниже берутся
         * исключительно из backend.
         */

        const image =
            document.createElement("div");

        image.className =
            "project-image";


        if (project.image_url) {

            image.style.backgroundImage =
                `url("${project.image_url}")`;

        }


        const content =
            document.createElement("div");

        content.className =
            "project-content";


        const top =
            document.createElement("div");

        top.className =
            "project-top";


        const title =
            document.createElement("h3");

        title.className =
            "project-title";

        title.textContent =
            project.title ?? "";


        const status =
            document.createElement("span");

        status.className =
            "project-status";

        status.textContent =
            project.status ?? "";


        top.appendChild(title);
        top.appendChild(status);


        const meta =
            document.createElement("div");

        meta.className =
            "project-meta";


        if (project.id) {

            const id =
                document.createElement("span");

            id.textContent =
                `№ ${project.id}`;

            meta.appendChild(id);
        }


        if (project.created_at) {

            const date =
                document.createElement("span");

            date.textContent =
                formatDate(project.created_at);

            meta.appendChild(date);
        }


        content.appendChild(top);
        content.appendChild(meta);


        article.appendChild(image);
        article.appendChild(content);


        if (project.url) {

            article.style.cursor =
                "pointer";

            article.addEventListener(
                "click",
                () => {
                    window.location.href =
                        project.url;
                }
            );
        }


        return article;
    }


    /*
     * ==========================================
     * FAVORITES
     * ==========================================
     */

    async function loadFavorites() {

        favoritesLoading.classList.remove("hidden");

        favoritesEmpty.classList.add("hidden");

        favoritesList.innerHTML = "";


        try {

            const favorites =
                await apiRequest(API.favorites);


            favoritesLoading.classList.add("hidden");


            if (
                !favorites ||
                !Array.isArray(favorites) ||
                favorites.length === 0
            ) {

                favoritesEmpty.classList.remove("hidden");

                return;
            }


            favorites.forEach(design => {

                favoritesList.appendChild(
                    createFavoriteCard(design)
                );

            });

        } catch (error) {

            favoritesLoading.classList.add("hidden");

            console.error(
                "Ошибка загрузки избранного:",
                error
            );

            favoritesEmpty.classList.remove("hidden");
        }
    }


    /*
     * ==========================================
     * FAVORITE CARD
     * ==========================================
     */

    function createFavoriteCard(design) {

        const article =
            document.createElement("article");

        article.className =
            "favorite-card";


        const image =
            document.createElement("div");

        image.className =
            "favorite-image";


        if (design.image_url) {

            image.style.backgroundImage =
                `url("${design.image_url}")`;

        }


        const content =
            document.createElement("div");

        content.className =
            "favorite-content";


        const title =
            document.createElement("h3");

        title.className =
            "favorite-title";

        title.textContent =
            design.title ?? "";


        const designer =
            document.createElement("p");

        designer.className =
            "favorite-designer";

        designer.textContent =
            design.designer_name ?? "";


        content.appendChild(title);
        content.appendChild(designer);


        article.appendChild(image);
        article.appendChild(content);


        if (design.url) {

            article.style.cursor =
                "pointer";

            article.addEventListener(
                "click",
                () => {
                    window.location.href =
                        design.url;
                }
            );
        }


        return article;
    }


    /*
     * ==========================================
     * PROFILE UPDATE
     * ==========================================
     */

    profileForm.addEventListener(
        "submit",
        async event => {

            event.preventDefault();


            const message =
                document.getElementById(
                    "profileMessage"
                );


            const name =
                document.getElementById(
                    "name"
                ).value.trim();


            if (!name) {

                showMessage(
                    message,
                    "Введите имя.",
                    true
                );

                return;
            }


            try {

                await apiRequest(
                    API.updateProfile,
                    {
                        method: "PUT",

                        body: JSON.stringify({
                            name: name
                        })
                    }
                );


                showMessage(
                    message,
                    "Изменения сохранены."
                );

            } catch (error) {

                showMessage(
                    message,
                    error.message,
                    true
                );
            }
        }
    );


    /*
     * ==========================================
     * PASSWORD
     * ==========================================
     */

    passwordForm.addEventListener(
        "submit",
        async event => {

            event.preventDefault();


            const message =
                document.getElementById(
                    "passwordMessage"
                );


            const currentPassword =
                document.getElementById(
                    "currentPassword"
                ).value;


            const newPassword =
                document.getElementById(
                    "newPassword"
                ).value;


            const confirmPassword =
                document.getElementById(
                    "confirmPassword"
                ).value;


            if (!currentPassword || !newPassword) {

                showMessage(
                    message,
                    "Заполните все необходимые поля.",
                    true
                );

                return;
            }


            if (
                newPassword !==
                confirmPassword
            ) {

                showMessage(
                    message,
                    "Пароли не совпадают.",
                    true
                );

                return;
            }


            try {

                await apiRequest(
                    API.changePassword,
                    {
                        method: "PUT",

                        body: JSON.stringify({
                            current_password:
                                currentPassword,

                            new_password:
                                newPassword
                        })
                    }
                );


                passwordForm.reset();


                showMessage(
                    message,
                    "Пароль успешно изменён."
                );

            } catch (error) {

                showMessage(
                    message,
                    error.message,
                    true
                );
            }
        }
    );


    /*
     * ==========================================
     * LOGOUT
     * ==========================================
     */

    logoutButton.addEventListener(
        "click",
        async () => {

            try {

                await apiRequest(
                    API.logout,
                    {
                        method: "POST"
                    }
                );

            } catch (error) {

                console.error(
                    "Ошибка выхода:",
                    error
                );

            } finally {

                window.location.href =
                    "/login";
            }
        }
    );


    /*
     * ==========================================
     * NAVIGATION
     * ==========================================
     */

    navItems.forEach(item => {

        item.addEventListener(
            "click",
            () => {

                const section =
                    item.dataset.section;


                if (!section) {
                    return;
                }


                navItems.forEach(
                    navItem => {
                        navItem.classList.remove(
                            "active"
                        );
                    }
                );


                panels.forEach(
                    panel => {
                        panel.classList.remove(
                            "active"
                        );
                    }
                );


                item.classList.add("active");


                const target =
                    document.getElementById(
                        `${section}Section`
                    );


                if (target) {

                    target.classList.add(
                        "active"
                    );
                }


                if (section === "projects") {

                    loadProjects();

                }


                if (section === "favorites") {

                    loadFavorites();

                }

            }
        );

    });


    /*
     * ==========================================
     * MESSAGE
     * ==========================================
     */

    function showMessage(
        element,
        text,
        error = false
    ) {

        element.textContent =
            text;

        element.classList.remove(
            "hidden",
            "error",
            "success"
        );


        element.classList.add(
            error
                ? "error"
                : "success"
        );


        setTimeout(
            () => {

                element.classList.add(
                    "hidden"
                );

            },
            4000
        );
    }


    /*
     * ==========================================
     * DATE
     * ==========================================
     */

    function formatDate(value) {

        const date =
            new Date(value);


        if (
            Number.isNaN(
                date.getTime()
            )
        ) {
            return "";
        }


        return date.toLocaleDateString(
            "ru-RU",
            {
                day: "2-digit",
                month: "2-digit",
                year: "numeric"
            }
        );
    }


    /*
     * ==========================================
     * MOBILE MENU
     * ==========================================
     */

    const mobileMenuButton =
        document.getElementById(
            "mobileMenuButton"
        );

    const mobileMenu =
        document.getElementById(
            "mobileMenu"
        );


    if (
        mobileMenuButton &&
        mobileMenu
    ) {

        mobileMenuButton.addEventListener(
            "click",
            () => {

                mobileMenu.classList.toggle(
                    "active"
                );

                mobileMenuButton.classList.toggle(
                    "active"
                );

            }
        );

    }


    /*
     * ==========================================
     * INITIALIZATION
     * ==========================================
     */

    document.addEventListener(
        "DOMContentLoaded",
        async () => {

            await loadProfile();

            // await loadProjects();

        }
    );
