def load_nav_obj(request):
    pages = [
            {
                "name": "home",
                "url": "home",
                "urls": ["home"],
                },
            {
                "name": "posts",
                "url": "posts",
                "urls": ["post", "posts"],
                }
            ]
    if request.user.is_authenticated:
        pages.append({
            "name": "nextcloud",
            "url": "https://nextcloud.aino-spring.com",
            "urls": []
            })
    return {"NAV_PAGES": pages}


def load_contact(request):
    return {
            "CONTACT": {
                "github": "https://github.com/theaino",
                "github_site": "https://github.com/theaino/aino_site",
                "instagram": "https://instagram.com/aino.spring",
                "email": "info@aino-spring.com",
                }
            }


def load_sites(request):
    return {
            "SITES": {
                "searxng": "https://search.aino-spring.com"
                }
            }
